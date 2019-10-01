package util

import (
	"encoding/xml"
)

type metadataJSON struct {
	Contact contact
}

type contact struct {
}

type metadataXML struct {
	XMLName        xml.Name `xml:"MD_Metadata"`
	Text           string   `xml:",chardata"`
	Gmd            string   `xml:"gmd,attr"`
	Geonet         string   `xml:"geonet,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	FileIdentifier struct {
		Text            string `xml:",chardata"`
		CharacterString struct {
			Text string `xml:",chardata"`
			Gco  string `xml:"gco,attr"`
		} `xml:"CharacterString"`
	} `xml:"fileIdentifier"`
	Language struct {
		Text         string `xml:",chardata"`
		LanguageCode struct {
			Text          string `xml:",chardata"`
			Xmlns         string `xml:"xmlns,attr"`
			CodeList      string `xml:"codeList,attr"`
			CodeListValue string `xml:"codeListValue,attr"`
		} `xml:"LanguageCode"`
	} `xml:"language"`
	CharacterSet struct {
		Text               string `xml:",chardata"`
		MDCharacterSetCode struct {
			Text          string `xml:",chardata"`
			Xmlns         string `xml:"xmlns,attr"`
			CodeList      string `xml:"codeList,attr"`
			CodeListValue string `xml:"codeListValue,attr"`
		} `xml:"MD_CharacterSetCode"`
	} `xml:"characterSet"`
	HierarchyLevel struct {
		Text        string `xml:",chardata"`
		MDScopeCode struct {
			Text          string `xml:",chardata"`
			Xmlns         string `xml:"xmlns,attr"`
			CodeList      string `xml:"codeList,attr"`
			CodeListValue string `xml:"codeListValue,attr"`
		} `xml:"MD_ScopeCode"`
	} `xml:"hierarchyLevel"`
	HierarchyLevelName struct {
		Text            string `xml:",chardata"`
		CharacterString struct {
			Text string `xml:",chardata"`
			Gco  string `xml:"gco,attr"`
		} `xml:"CharacterString"`
	} `xml:"hierarchyLevelName"`
	Contact struct {
		Text               string `xml:",chardata"`
		CIResponsibleParty struct {
			Text           string `xml:",chardata"`
			IndividualName struct {
				Text            string `xml:",chardata"`
				CharacterString struct {
					Text string `xml:",chardata"`
					Gco  string `xml:"gco,attr"`
				} `xml:"CharacterString"`
			} `xml:"individualName"`
			OrganisationName struct {
				Text            string `xml:",chardata"`
				CharacterString struct {
					Text string `xml:",chardata"`
					Gco  string `xml:"gco,attr"`
				} `xml:"CharacterString"`
			} `xml:"organisationName"`
			ContactInfo struct {
				Text      string `xml:",chardata"`
				CIContact struct {
					Text  string `xml:",chardata"`
					Phone struct {
						Text        string `xml:",chardata"`
						CITelephone struct {
							Text  string `xml:",chardata"`
							Voice struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"voice"`
							Facsimile struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"facsimile"`
						} `xml:"CI_Telephone"`
					} `xml:"phone"`
					Address struct {
						Text      string `xml:",chardata"`
						CIAddress struct {
							Text          string `xml:",chardata"`
							DeliveryPoint struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"deliveryPoint"`
							City struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"city"`
							AdministrativeArea struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"administrativeArea"`
							PostalCode struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"postalCode"`
							Country struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"country"`
							ElectronicMailAddress struct {
								Text            string `xml:",chardata"`
								CharacterString struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"CharacterString"`
							} `xml:"electronicMailAddress"`
						} `xml:"CI_Address"`
					} `xml:"address"`
				} `xml:"CI_Contact"`
			} `xml:"contactInfo"`
			Role struct {
				Text       string `xml:",chardata"`
				CIRoleCode struct {
					Text          string `xml:",chardata"`
					CodeList      string `xml:"codeList,attr"`
					CodeListValue string `xml:"codeListValue,attr"`
				} `xml:"CI_RoleCode"`
			} `xml:"role"`
		} `xml:"CI_ResponsibleParty"`
	} `xml:"contact"`
	DateStamp struct {
		Text string `xml:",chardata"`
		Date struct {
			Text string `xml:",chardata"`
			Gco  string `xml:"gco,attr"`
		} `xml:"Date"`
	} `xml:"dateStamp"`
	MetadataStandardName struct {
		Text            string `xml:",chardata"`
		CharacterString struct {
			Text string `xml:",chardata"`
			Gco  string `xml:"gco,attr"`
		} `xml:"CharacterString"`
	} `xml:"metadataStandardName"`
	MetadataStandardVersion struct {
		Text            string `xml:",chardata"`
		CharacterString struct {
			Text string `xml:",chardata"`
			Gco  string `xml:"gco,attr"`
		} `xml:"CharacterString"`
	} `xml:"metadataStandardVersion"`
	ReferenceSystemInfo struct {
		Text              string `xml:",chardata"`
		MDReferenceSystem struct {
			Text                      string `xml:",chardata"`
			ReferenceSystemIdentifier struct {
				Text         string `xml:",chardata"`
				RSIdentifier struct {
					Text string `xml:",chardata"`
					Code struct {
						Text            string `xml:",chardata"`
						CharacterString struct {
							Text string `xml:",chardata"`
							Gco  string `xml:"gco,attr"`
						} `xml:"CharacterString"`
					} `xml:"code"`
					CodeSpace struct {
						Text            string `xml:",chardata"`
						CharacterString struct {
							Text string `xml:",chardata"`
							Gco  string `xml:"gco,attr"`
						} `xml:"CharacterString"`
					} `xml:"codeSpace"`
					Version struct {
						Text            string `xml:",chardata"`
						CharacterString struct {
							Text string `xml:",chardata"`
							Gco  string `xml:"gco,attr"`
						} `xml:"CharacterString"`
					} `xml:"version"`
				} `xml:"RS_Identifier"`
			} `xml:"referenceSystemIdentifier"`
		} `xml:"MD_ReferenceSystem"`
	} `xml:"referenceSystemInfo"`
	IdentificationInfo struct {
		Text                 string `xml:",chardata"`
		Gco                  string `xml:"gco,attr"`
		MDDataIdentification struct {
			Text     string `xml:",chardata"`
			Citation struct {
				Text       string `xml:",chardata"`
				CICitation struct {
					Text  string `xml:",chardata"`
					Title struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"title"`
					Date []struct {
						Text   string `xml:",chardata"`
						CIDate struct {
							Text string `xml:",chardata"`
							Date struct {
								Text string `xml:",chardata"`
								Date string `xml:"Date"`
							} `xml:"date"`
							DateType struct {
								Text           string `xml:",chardata"`
								CIDateTypeCode struct {
									Text          string `xml:",chardata"`
									Xmlns         string `xml:"xmlns,attr"`
									CodeList      string `xml:"codeList,attr"`
									CodeListValue string `xml:"codeListValue,attr"`
								} `xml:"CI_DateTypeCode"`
							} `xml:"dateType"`
						} `xml:"CI_Date"`
					} `xml:"date"`
					Identifier struct {
						Text         string `xml:",chardata"`
						MDIdentifier struct {
							Text string `xml:",chardata"`
							Code struct {
								Text            string `xml:",chardata"`
								CharacterString string `xml:"CharacterString"`
							} `xml:"code"`
						} `xml:"MD_Identifier"`
					} `xml:"identifier"`
				} `xml:"CI_Citation"`
			} `xml:"citation"`
			Abstract struct {
				Text            string `xml:",chardata"`
				CharacterString string `xml:"CharacterString"`
			} `xml:"abstract"`
			Status struct {
				Text           string `xml:",chardata"`
				MDProgressCode struct {
					Text          string `xml:",chardata"`
					Xmlns         string `xml:"xmlns,attr"`
					CodeList      string `xml:"codeList,attr"`
					CodeListValue string `xml:"codeListValue,attr"`
				} `xml:"MD_ProgressCode"`
			} `xml:"status"`
			PointOfContact struct {
				Text               string `xml:",chardata"`
				CIResponsibleParty struct {
					Text           string `xml:",chardata"`
					IndividualName struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"individualName"`
					OrganisationName struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"organisationName"`
					ContactInfo struct {
						Text      string `xml:",chardata"`
						CIContact struct {
							Text  string `xml:",chardata"`
							Phone struct {
								Text        string `xml:",chardata"`
								CITelephone struct {
									Text  string `xml:",chardata"`
									Voice struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"voice"`
									Facsimile struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"facsimile"`
								} `xml:"CI_Telephone"`
							} `xml:"phone"`
							Address struct {
								Text      string `xml:",chardata"`
								CIAddress struct {
									Text          string `xml:",chardata"`
									DeliveryPoint struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"deliveryPoint"`
									City struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"city"`
									AdministrativeArea struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"administrativeArea"`
									PostalCode struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"postalCode"`
									Country struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"country"`
									ElectronicMailAddress struct {
										Text            string `xml:",chardata"`
										CharacterString string `xml:"CharacterString"`
									} `xml:"electronicMailAddress"`
								} `xml:"CI_Address"`
							} `xml:"address"`
						} `xml:"CI_Contact"`
					} `xml:"contactInfo"`
					Role struct {
						Text       string `xml:",chardata"`
						CIRoleCode struct {
							Text          string `xml:",chardata"`
							CodeList      string `xml:"codeList,attr"`
							CodeListValue string `xml:"codeListValue,attr"`
						} `xml:"CI_RoleCode"`
					} `xml:"role"`
				} `xml:"CI_ResponsibleParty"`
			} `xml:"pointOfContact"`
			ResourceMaintenance struct {
				Text                     string `xml:",chardata"`
				MDMaintenanceInformation struct {
					Text                          string `xml:",chardata"`
					MaintenanceAndUpdateFrequency struct {
						Text                       string `xml:",chardata"`
						MDMaintenanceFrequencyCode struct {
							Text          string `xml:",chardata"`
							Xmlns         string `xml:"xmlns,attr"`
							CodeList      string `xml:"codeList,attr"`
							CodeListValue string `xml:"codeListValue,attr"`
						} `xml:"MD_MaintenanceFrequencyCode"`
					} `xml:"maintenanceAndUpdateFrequency"`
					MaintenanceNote struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"maintenanceNote"`
				} `xml:"MD_MaintenanceInformation"`
			} `xml:"resourceMaintenance"`
			GraphicOverview struct {
				Text            string `xml:",chardata"`
				MDBrowseGraphic struct {
					Text     string `xml:",chardata"`
					FileName struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"fileName"`
				} `xml:"MD_BrowseGraphic"`
			} `xml:"graphicOverview"`
			DescriptiveKeywords []struct {
				Text       string `xml:",chardata"`
				MDKeywords struct {
					Text    string `xml:",chardata"`
					Keyword []struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"keyword"`
					ThesaurusName struct {
						Text       string `xml:",chardata"`
						CICitation struct {
							Text  string `xml:",chardata"`
							Title struct {
								Text            string `xml:",chardata"`
								CharacterString string `xml:"CharacterString"`
							} `xml:"title"`
							Date struct {
								Text   string `xml:",chardata"`
								CIDate struct {
									Text string `xml:",chardata"`
									Date struct {
										Text string `xml:",chardata"`
										Date string `xml:"Date"`
									} `xml:"date"`
									DateType struct {
										Text           string `xml:",chardata"`
										CIDateTypeCode struct {
											Text          string `xml:",chardata"`
											Xmlns         string `xml:"xmlns,attr"`
											CodeList      string `xml:"codeList,attr"`
											CodeListValue string `xml:"codeListValue,attr"`
										} `xml:"CI_DateTypeCode"`
									} `xml:"dateType"`
								} `xml:"CI_Date"`
							} `xml:"date"`
						} `xml:"CI_Citation"`
					} `xml:"thesaurusName"`
				} `xml:"MD_Keywords"`
			} `xml:"descriptiveKeywords"`
			ResourceConstraints []struct {
				Text               string `xml:",chardata"`
				MDLegalConstraints struct {
					Text          string `xml:",chardata"`
					UseLimitation struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"useLimitation"`
					AccessConstraints struct {
						Text              string `xml:",chardata"`
						MDRestrictionCode struct {
							Text          string `xml:",chardata"`
							Xmlns         string `xml:"xmlns,attr"`
							CodeList      string `xml:"codeList,attr"`
							CodeListValue string `xml:"codeListValue,attr"`
						} `xml:"MD_RestrictionCode"`
					} `xml:"accessConstraints"`
					OtherConstraints struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"otherConstraints"`
					UseConstraints struct {
						Text              string `xml:",chardata"`
						MDRestrictionCode struct {
							Text          string `xml:",chardata"`
							Xmlns         string `xml:"xmlns,attr"`
							CodeList      string `xml:"codeList,attr"`
							CodeListValue string `xml:"codeListValue,attr"`
						} `xml:"MD_RestrictionCode"`
					} `xml:"useConstraints"`
				} `xml:"MD_LegalConstraints"`
			} `xml:"resourceConstraints"`
			SpatialResolution struct {
				Text         string `xml:",chardata"`
				MDResolution struct {
					Text            string `xml:",chardata"`
					EquivalentScale struct {
						Text                     string `xml:",chardata"`
						MDRepresentativeFraction struct {
							Text        string `xml:",chardata"`
							Denominator struct {
								Text    string `xml:",chardata"`
								Integer string `xml:"Integer"`
							} `xml:"denominator"`
						} `xml:"MD_RepresentativeFraction"`
					} `xml:"equivalentScale"`
				} `xml:"MD_Resolution"`
			} `xml:"spatialResolution"`
			Language struct {
				Text         string `xml:",chardata"`
				LanguageCode struct {
					Text          string `xml:",chardata"`
					Xmlns         string `xml:"xmlns,attr"`
					CodeList      string `xml:"codeList,attr"`
					CodeListValue string `xml:"codeListValue,attr"`
				} `xml:"LanguageCode"`
			} `xml:"language"`
			TopicCategory struct {
				Text                string `xml:",chardata"`
				MDTopicCategoryCode string `xml:"MD_TopicCategoryCode"`
			} `xml:"topicCategory"`
			Extent struct {
				Text     string `xml:",chardata"`
				EXExtent struct {
					Text              string `xml:",chardata"`
					GeographicElement struct {
						Text                    string `xml:",chardata"`
						EXGeographicBoundingBox struct {
							Text           string `xml:",chardata"`
							ExtentTypeCode struct {
								Text    string `xml:",chardata"`
								Boolean string `xml:"Boolean"`
							} `xml:"extentTypeCode"`
							WestBoundLongitude struct {
								Text    string `xml:",chardata"`
								Decimal string `xml:"Decimal"`
							} `xml:"westBoundLongitude"`
							EastBoundLongitude struct {
								Text    string `xml:",chardata"`
								Decimal string `xml:"Decimal"`
							} `xml:"eastBoundLongitude"`
							SouthBoundLatitude struct {
								Text    string `xml:",chardata"`
								Decimal string `xml:"Decimal"`
							} `xml:"southBoundLatitude"`
							NorthBoundLatitude struct {
								Text    string `xml:",chardata"`
								Decimal string `xml:"Decimal"`
							} `xml:"northBoundLatitude"`
						} `xml:"EX_GeographicBoundingBox"`
					} `xml:"geographicElement"`
				} `xml:"EX_Extent"`
			} `xml:"extent"`
		} `xml:"MD_DataIdentification"`
	} `xml:"identificationInfo"`
	DistributionInfo struct {
		Text           string `xml:",chardata"`
		Gco            string `xml:"gco,attr"`
		MDDistribution struct {
			Text               string `xml:",chardata"`
			DistributionFormat []struct {
				Text     string `xml:",chardata"`
				MDFormat struct {
					Text string `xml:",chardata"`
					Name struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"name"`
					Version struct {
						Text            string `xml:",chardata"`
						CharacterString string `xml:"CharacterString"`
					} `xml:"version"`
				} `xml:"MD_Format"`
			} `xml:"distributionFormat"`
			TransferOptions []struct {
				Text                     string `xml:",chardata"`
				MDDigitalTransferOptions struct {
					Text   string `xml:",chardata"`
					OnLine struct {
						Text             string `xml:",chardata"`
						CIOnlineResource struct {
							Text    string `xml:",chardata"`
							Linkage struct {
								Text string `xml:",chardata"`
								URL  string `xml:"URL"`
							} `xml:"linkage"`
							Function struct {
								Text                 string `xml:",chardata"`
								CIOnLineFunctionCode struct {
									Text          string `xml:",chardata"`
									Xmlns         string `xml:"xmlns,attr"`
									CodeList      string `xml:"codeList,attr"`
									CodeListValue string `xml:"codeListValue,attr"`
								} `xml:"CI_OnLineFunctionCode"`
							} `xml:"function"`
						} `xml:"CI_OnlineResource"`
					} `xml:"onLine"`
				} `xml:"MD_DigitalTransferOptions"`
			} `xml:"transferOptions"`
		} `xml:"MD_Distribution"`
	} `xml:"distributionInfo"`
	DataQualityInfo struct {
		Text          string `xml:",chardata"`
		DQDataQuality struct {
			Text  string `xml:",chardata"`
			Scope struct {
				Text    string `xml:",chardata"`
				DQScope struct {
					Text  string `xml:",chardata"`
					Level struct {
						Text        string `xml:",chardata"`
						MDScopeCode struct {
							Text          string `xml:",chardata"`
							Xmlns         string `xml:"xmlns,attr"`
							CodeList      string `xml:"codeList,attr"`
							CodeListValue string `xml:"codeListValue,attr"`
						} `xml:"MD_ScopeCode"`
					} `xml:"level"`
					LevelDescription struct {
						Text      string `xml:",chardata"`
						Gco       string `xml:"gco,attr"`
						NilReason string `xml:"nilReason,attr"`
					} `xml:"levelDescription"`
				} `xml:"DQ_Scope"`
			} `xml:"scope"`
			Report struct {
				Text                string `xml:",chardata"`
				DQDomainConsistency struct {
					Text   string `xml:",chardata"`
					Result struct {
						Text                string `xml:",chardata"`
						DQConformanceResult struct {
							Text          string `xml:",chardata"`
							Specification struct {
								Text       string `xml:",chardata"`
								CICitation struct {
									Text  string `xml:",chardata"`
									Title struct {
										Text            string `xml:",chardata"`
										CharacterString struct {
											Text string `xml:",chardata"`
											Gco  string `xml:"gco,attr"`
										} `xml:"CharacterString"`
									} `xml:"title"`
									Date struct {
										Text   string `xml:",chardata"`
										CIDate struct {
											Text string `xml:",chardata"`
											Date struct {
												Text string `xml:",chardata"`
												Date struct {
													Text string `xml:",chardata"`
													Gco  string `xml:"gco,attr"`
												} `xml:"Date"`
											} `xml:"date"`
											DateType struct {
												Text           string `xml:",chardata"`
												CIDateTypeCode struct {
													Text          string `xml:",chardata"`
													Xmlns         string `xml:"xmlns,attr"`
													CodeList      string `xml:"codeList,attr"`
													CodeListValue string `xml:"codeListValue,attr"`
												} `xml:"CI_DateTypeCode"`
											} `xml:"dateType"`
										} `xml:"CI_Date"`
									} `xml:"date"`
								} `xml:"CI_Citation"`
							} `xml:"specification"`
							Explanation struct {
								Text      string `xml:",chardata"`
								Gco       string `xml:"gco,attr"`
								NilReason string `xml:"nilReason,attr"`
							} `xml:"explanation"`
							Pass struct {
								Text    string `xml:",chardata"`
								Boolean struct {
									Text string `xml:",chardata"`
									Gco  string `xml:"gco,attr"`
								} `xml:"Boolean"`
							} `xml:"pass"`
						} `xml:"DQ_ConformanceResult"`
					} `xml:"result"`
				} `xml:"DQ_DomainConsistency"`
			} `xml:"report"`
			Lineage struct {
				Text      string `xml:",chardata"`
				LILineage struct {
					Text      string `xml:",chardata"`
					Statement struct {
						Text            string `xml:",chardata"`
						CharacterString struct {
							Text string `xml:",chardata"`
							Gco  string `xml:"gco,attr"`
						} `xml:"CharacterString"`
					} `xml:"statement"`
				} `xml:"LI_Lineage"`
			} `xml:"lineage"`
		} `xml:"DQ_DataQuality"`
	} `xml:"dataQualityInfo"`
}

// ParseMetadataToJSON converts ISO 19115, as well as INSPIRE verified metadata from XML to JSON
func ParseMetadataToJSON(raw []byte) []byte {
	var mdXML metadataXML

	err := xml.Unmarshal(raw, &mdXML)
	if err != nil {
		return nil
	}

	//mdJSON := metadataJSON{}

	return nil
}
